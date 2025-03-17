#!/bin/sh
set -euo pipefail

export AWS_PAGER=""

KEYS="sweep selena bella husky german_sheperd mr_peanutbutter luna"

for KEY in $KEYS; do
  aws rekognition detect-labels \
    --image '{"S3Object":{"Bucket":"rekog-sandbox","Name":"'"${KEY}"'.jpeg"}}' \
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
    }' > $KEY.json
done



