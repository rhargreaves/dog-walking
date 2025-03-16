#!/bin/sh
export AWS_PAGER=""

aws rekognition detect-labels \
  --image '{"S3Object":{"Bucket":"uat-dog-images","Name":"1896a038-473a-4c89-b70a-ebe465045dee"}}' \
  --min-confidence 55 \
  --max-labels 10 \
  --region eu-west-1 \
  --features 'GENERAL_LABELS' \
  --settings '{
    "GeneralLabels": {
      "LabelCategoryInclusionFilters": [
        "Animals and Pets"
      ],
      "LabelExclusionFilters": [
        "Pet",
        "Mammal",
        "Canine",
        "Animal"
      ]
    }
  }'
